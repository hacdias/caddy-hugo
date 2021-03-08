
const ThreeDModelExtRegExp = new RegExp(/\.(obj|stl|dae|ply|fbx|gltf)$/);

export const is3DModelFile = (filename) => {
    return ThreeDModelExtRegExp.test(filename)
}