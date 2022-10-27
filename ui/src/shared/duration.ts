import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

function fromNow(ts: string) {
  const from = dayjs(ts)

  return from.fromNow()
}

export {fromNow}
