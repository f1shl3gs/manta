describe('handleMove', () => {
  test('swap', () => {
    const arr = ['a', 'b']
    const oldIndex = 0
    const newIndex = 1

    const tmp = arr[oldIndex]
    arr[oldIndex] = arr[newIndex]
    arr[newIndex] = tmp

    console.log('arr', arr)
  })
})
