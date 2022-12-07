export type Result = {
  metric: {
    [key: string]: string
  }
  values: [[number, string]]
}

export interface Row {
  [key: string]: string | number
}

export type PromResp = {
  status: string
  data: {
    resultType: string
    result: Result[]
  }
}

export const transformToRows = (resp: PromResp): Row[] => {
  if (!resp) {
    return []
  }

  if (resp.status !== 'success') {
    return []
  }

  return resp.data.result
    .map((item: Result) => {
      const {metric, values} = item

      return values.map(val => {
        return {
          ...metric,
          _time: val[0] * 1000,
          _value: Number(val[1]),
        }
      })
    })
    .flat()
}
