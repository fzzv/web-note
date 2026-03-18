function twoSum(nums: number[], target: number): number[] | null {
  for (let i = 0; i < nums.length; i++) {
    for (let j = 0; j < nums.length; j++) {
      if (nums[i] + nums[j] === target && i !== j) {
        return [i, j];
      }
    }
  }
  return null;
}

function twoSumHash(nums: number[], target: number): number[] | null {
  const map = new Map<number, number>();
  for (let i = 0; i < nums.length; i++) {
    const need = target - nums[i];
    if (map.has(need)) {
      return [map.get(need) as number, i];
    }
    map.set(nums[i], i);
  }
  return null;
}
