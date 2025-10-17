export const cloneDeep = <T>(obj: T): T => {
  return JSON.parse(JSON.stringify(obj))
}
