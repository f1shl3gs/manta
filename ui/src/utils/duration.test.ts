import {parseDuration} from './duration'

describe('duration', () => {
  test('parse simple duration', () => {
    const result = parseDuration('15s')
    console.log(result)
  })

  test('parse multi', () => {
    const result = parseDuration('20m15s')
    console.log(result)
  })
})
