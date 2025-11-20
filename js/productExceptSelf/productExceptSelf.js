/**
 * @param {number[]} nums
 * @return {number[]}
 */
const productExceptSelf = (nums) => {
    const n = nums.length;
    const result = [];
    
    let leftProduct = 1;
    let rightProduct = 1;
    
    for (let i = 0; i < n; i++) {
        result[i] = leftProduct;
        leftProduct = leftProduct * nums[i];
    }

    for (let i = n - 1; i >= 0; i--) {
        result[i] = result[i] * rightProduct;
        rightProduct = rightProduct * nums[i];
    } 

    return result;
};

const nums = [1,2,3,4]

console.log(productExceptSelf(nums))
