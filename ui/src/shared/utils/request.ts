export interface RequestOptions {
  method?: string
  body?: object
  query?: string[][] | Record<string, string> | string | URLSearchParams
}

const request = async function (
  url: string,
  options?: RequestOptions
): Promise<any> {
  const body = options?.body ? JSON.stringify(options.body) : null
  const query = options?.query ? `?${new URLSearchParams(options.query)}` : ''

  const resp = await fetch(`${url}${query}`, {
    method: options?.method || 'GET',
    body,
  })
  const {status, headers} = resp

  const respContentType = headers.get('Content-Type') || ''

  let data
  if (respContentType.includes('json')) {
    data = await resp.json()
  } else if (respContentType.includes('octet-stream')) {
    data = await resp.blob()
  } else {
    data = await resp.text()
  }

  return {
    status,
    headers,
    data,
  }
}

export default request
