export function swap(list: any[], i: number) : any[] {
  if (i <= 0) {
    return list
  }
  if (i >= list.length) {
    return list
  }
  [list[i], list[i - 1]] = [list[i - 1], list[i]]
  return list
}
