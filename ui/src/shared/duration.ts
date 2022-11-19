import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export const fromNow = (ts: string) => {
  const from = dayjs(ts)

  return from.fromNow()
}

export const parseDuration = (input: string): number => {
  const durationRegExp = /([0-9]+)(y|mo|w|d|h|ms|s|m|us|µs|ns)/g

  // warning! regex.exec(string) modifies the regex it is operating on so that subsequent calls on the same string behave differently
  let match = durationRegExp.exec(input)

  if (!match) {
    throw new Error(`could not parse "${input}" as duration`)
  }

  let d = 0
  while (match) {
    // ms
    let factor = 1
    switch (match[2]) {
      case 'ms':
        break
      case 's':
        factor = 1000
        break
      case 'm':
        factor = 60 * 1000
        break
      case 'h':
        factor = 60 * 60 * 1000
        break
      case 'd':
        factor = 24 * 60 * 60 * 1000
        break
      case 'w':
        factor = 7 * 24 * 60 * 60 * 1000
        break
      default:
        // eslint-disable-next-line no-throw-literal
        throw `Unsupported suffix "${match[2]}" `
    }

    d += Number(match[1]) * factor
    match = durationRegExp.exec(input)
  }

  return d
}
