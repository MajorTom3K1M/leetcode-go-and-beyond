const merge = (intervals) => {
    intervals.sort((a, b) => {
        return a[0] - b[0]
    })

    const merge = [intervals[0]];
    console.log(intervals)
    for (let i = 1; i < intervals.length; i++) {
        if (merge.length > 0 && merge[merge.length - 1][1] >= intervals[i][0]) {
            if (merge[merge.length - 1][1] < intervals[i][1]) {
                merge[merge.length - 1][1] = intervals[i][1]
            }             
        } else {
            merge.push(intervals[i])
        }
    }

    return merge
};

// const intervals = [[1,3],[2,6],[8,10],[15,18]]
// const intervals = [[1,3]]
// const intervals = [[1, 3], [2, 6], [8, 10], [15, 18], [16, 19]]
// const intervals = [[1,4],[0,4]]
const intervals = [[2,3],[4,5],[6,7],[8,9],[1,10]]

console.log(merge(intervals));