<template>
  <div id="editor-container">
    <header-bar>
      <action icon="close" :label="$t('buttons.close')" @action="close()" />
      <title>{{ req.name }}</title>
    </header-bar>

    <breadcrumbs base="/files" noLink />

    <div id="editor"></div>
  </div>
</template>

<style scoped>
#editor-container {
  height: 100vh;
  width: 100vw;
}
</style>

<script>
import { mapState } from "vuex";
import url from "@/utils/url";
import { baseURL, onlyOffice } from "@/utils/constants";

import HeaderBar from "@/components/header/HeaderBar";
import Action from "@/components/header/Action";
import Breadcrumbs from "@/components/Breadcrumbs";

export default {
  name: "onlyofficeeditor",
  components: {
    HeaderBar,
    Action,
    Breadcrumbs,
  },
  data: function () {
    return {};
  },
  computed: {
    ...mapState(["req", "user", "jwt"]),
    breadcrumbs() {
      let parts = this.$route.path.split("/");

      if (parts[0] === "") {
        parts.shift();
      }

      if (parts[parts.length - 1] === "") {
        parts.pop();
      }

      let breadcrumbs = [];

      for (let i = 0; i < parts.length; i++) {
        breadcrumbs.push({ name: decodeURIComponent(parts[i]) });
      }

      breadcrumbs.shift();

      if (breadcrumbs.length > 3) {
        while (breadcrumbs.length !== 4) {
          breadcrumbs.shift();
        }

        breadcrumbs[0].name = "...";
      }

      return breadcrumbs;
    },
  },
  created() {
    window.addEventListener("keydown", this.keyEvent);
  },
  beforeDestroy() {
    window.removeEventListener("keydown", this.keyEvent);
    this.editor.destroyEditor();
  },
  mounted: function () {
    let onlyofficeScript = document.createElement("script");
    onlyofficeScript.setAttribute(
      "src",
      `${onlyOffice}/web-apps/apps/api/documents/api.js`
    );
    document.head.appendChild(onlyofficeScript);

    onlyofficeScript.onload = () => {
      let fileUrl = `${window.location.protocol}//${window.location.host}${baseURL}/api/raw${url.encodePath(
        this.req.path
      )}?auth=${this.jwt}`;

      let key = Date.parse(this.req.modified).toString() + url.encodePath(this.req.path);
      key = key.replaceAll(/[-_.!~[\]*'()/,;:\-%+.]/g, "");
      if (key.length > 127) {
        key = key.substring(0, 127);
      }

      /*eslint-disable */
      let config = {
        document: {
          fileType: this.req.extension.substring(1),
          key: key,
          title: this.req.name,
          url: fileUrl,
          permissions: {
            edit: this.user.perm.modify,
            download: this.user.perm.download,
            print: this.user.perm.download
          }
        },
        editorConfig: {
          callbackUrl: `${window.location.protocol}//${window.location.host}${baseURL}/api/onlyoffice/callback?auth=${this.jwt}&save=${encodeURIComponent(this.req.path)}`,
          user: {
            id: this.user.id,
            name: `User ${this.user.id}`
          },
          customization: {
            autosave: true,
            forcesave: true
          },
          lang: this.user.locale,
          mode: this.user.perm.modify ? "edit" : "view"
        }
      };
      this.editor = new DocsAPI.DocEditor("editor", config);
      /*eslint-enable */
    };
  },
  methods: {
    back() {
      let uri = url.removeLastDir(this.$route.path) + "/";
      this.$router.push({ path: uri });
    },
    keyEvent(event) {
      if (!event.ctrlKey && !event.metaKey) {
        return;
      }

      if (String.fromCharCode(event.which).toLowerCase() !== "s") {
        return;
      }

      event.preventDefault();
      this.save();
    },
    close() {
      this.$store.commit("updateRequest", {});

      let uri = url.removeLastDir(this.$route.path) + "/";
      this.$router.push({ path: uri });
    },
  },
};
</script>
