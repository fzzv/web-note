export function mapToIds(values: any[]) {
  return JSON.stringify(values.map(item => item.id));
}
