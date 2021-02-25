import constate from 'constate'
import {useLocation} from 'react-router-dom'

function parseQuery(qs: string) {
  let query = {} as {[key: string]: string}
  let pairs = (qs[0] === '?' ? qs.substr(1) : qs).split('&')
  for (let i = 0; i < pairs.length; i++) {
    let pair = pairs[i].split('=')
    let key = decodeURIComponent(pair[0])
    query[key] = decodeURIComponent(pair[1] || '')
  }

  return query
}

const [SearchProvider, useSearch] = constate(
  () => {
    const location = useLocation()
    const query = parseQuery(location.search)

    return {
      ...query,
    }
  },
  value => value
)

export {SearchProvider, useSearch}
