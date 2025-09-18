export function def(value: any) {
  return value ? '"' + value + '"' : '""';
};
