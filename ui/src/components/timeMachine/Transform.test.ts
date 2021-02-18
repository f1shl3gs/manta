import {fromRows} from '@influxdata/giraffe'

const s1 = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {
          foo: 'bar',
          instance: '127.0.0.1:8080',
          job: 'selfstat',
        },
        values: [
          [1613377622, '1.5808588888888888'],
          [1613377636, '1.4076944444444444'],
          [1613377650, '1.111111111111111'],
          [1613377664, '0.04444444444444449'],
        ],
      },
    ],
  },
}

const s2 = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {
          __name__: 'prometheus_tsdb_head_series',
          foo: 'bar',
          instance: '127.0.0.1:8080',
          job: 'selfstat',
          tenant: '06fe69936c23d000',
        },
        values: [
          [1613377608, '686'],
          [1613377622, '363'],
          [1613377636, '460'],
          [1613377650, '460'],
          [1613377664, '460'],
        ],
      },
    ],
  },
}

describe('multi resp', () => {
  it('merge', () => {
    let resps = [s1, s2]
    let rows = []
    resps.forEach(resp => {
      resp.data.result.map(series => {
        series.values.forEach(sample => {
          rows.push({
            ...series.metric,
            time: sample[0],
            value: sample[1],
          })
        })
      })
    })

    const table = fromRows(rows)
    expect(table.length).toEqual(9)
  })
})
