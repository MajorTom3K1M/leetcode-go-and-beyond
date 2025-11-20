/**
 * @param {character[][]} board
 * @return {boolean}
 */
const isValidSudoku = (board) => {
    const duplicateMap = new Map();

    for (let col = 0; col < board.length; col++) {
        for (let row = 0; row < board[col].length; row++) {
            let val = board[row][col];
            if (!duplicateMap.has(val)) {
                duplicateMap.set(val, true);
            } else if (duplicateMap.has(val) && val != ".") {
                return false;
            }
        }
        duplicateMap.clear();
    }


    for (let col = 0; col < board.length; col++) {
        for (let row = 0; row < board[col].length; row++) {
            let val = board[col][row];
            if (!duplicateMap.has(val)) {
                duplicateMap.set(val, true);
            } else if (duplicateMap.has(val) && val != ".") {
                return false
            }
        }
        duplicateMap.clear();
    }

    for (let i = 0; i < board.length; i++) {
        for (let j = 0; j < board[i].length; j++) {
            let col = Math.floor(i / 3)
            let row = Math.floor(j / 3)
            let isHasNestedMap = duplicateMap.has(`${row},${col}`);
            let val = board[i][j];
            if (!isHasNestedMap) {
                duplicateMap.set(`${row},${col}`, new Map());
            }

            let nestedMap = duplicateMap.get(`${row},${col}`);

            if (!nestedMap.has(val)) {
                nestedMap.set(val, true);
            } else if (nestedMap.has(val) && val != ".") {
                return false;
            }
        }
    }

    return true;
};

board = [["5", "3", ".", ".", "7", ".", ".", ".", "."]
    , ["6", ".", ".", "1", "9", "5", ".", ".", "."]
    , [".", "9", "8", ".", ".", ".", ".", "6", "."]
    , ["8", ".", ".", ".", "6", ".", ".", ".", "3"]
    , ["4", ".", ".", "8", ".", "3", ".", ".", "1"]
    , ["7", ".", ".", ".", "2", ".", ".", ".", "6"]
    , [".", "6", ".", ".", ".", ".", "2", "8", "."]
    , [".", ".", ".", "4", "1", "9", ".", ".", "5"]
    , [".", ".", ".", ".", "8", ".", ".", "7", "9"]]

console.log(isValidSudoku(board))