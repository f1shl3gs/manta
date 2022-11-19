import {FromFluxResult, fromRows} from '@influxdata/giraffe'

export type Result = {
  metric: {
    [key: string]: string
  }
  values: [[number, string]]
}

export type PromResp = {
  status: string
  data: {
    resultType: string
    result: Result[]
  }
}

export interface Row {
  [key: string]: string | number
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
          time: val[0] * 1000,
          value: Number(val[1]),
        }
      })
    })
    .flat()
}

export const transformPromResp = (
  resp?: PromResp
): FromFluxResult | undefined => {
  if (!resp) {
    return undefined
  }

  if (resp.status !== 'success') {
    return undefined
  }

  const rows = resp.data.result
    .map((item: Result) => {
      const {metric, values} = item

      return values.map(val => {
        return {
          ...metric,
          time: val[0] * 1000,
          value: Number(val[1]),
        }
      })
    })
    .flat()

  const table = fromRows(rows)
  const groupKeys = table.columnKeys.filter(
    key => key !== 'time' && key !== 'value'
  )

  return {
    table,
    fluxGroupKeyUnion: groupKeys,
    // TODO:
    resultColumnNames: [],
  }
}
