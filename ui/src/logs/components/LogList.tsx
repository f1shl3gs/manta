// libraries
import React, {useMemo} from 'react'
import moment from 'moment'

// components
import {ComponentSize, Table} from '@influxdata/clockface'

// test data
import testData from './query_range_resp.json'

type StreamResult = {
  stream: {
    [key: string]: string
  }
  values: string[][]
}

type MatrixResult = {
  metric: {
    [key: string]: string
  }
  values: [[number, string]]
}

type Resp = {
  status: string
  data: {
    resultType: string
    result: StreamResult[] | MatrixResult[]
  }
}

type TableRow = {
  ts: string
  msg: string
}

const transformData = (ss: StreamResult[]) => {
  const rows = new Array<TableRow>()

  ss.forEach((stream) => {
    stream.values.forEach((pair) => {
      rows.push({
        ts: moment(Number(pair[0]) / 1000 / 1000).format(),
        msg: pair[1],
      })
    })
  })

  return rows
}

const LogList = () => {
  // todo: show the common labels

  const resp = testData as Resp
  const rows = useMemo(() => {
    return transformData(resp.data.result as StreamResult[])
  }, [resp.data.result])

  return (
    <Table fontSize={ComponentSize.ExtraSmall}>
      <Table.Header>
        <Table.Row>
          <Table.HeaderCell>Timestamp</Table.HeaderCell>
          <Table.HeaderCell>Message</Table.HeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {rows.map((item, index) => (
          <Table.Row key={index}>
            <Table.Cell>{item.ts}</Table.Cell>
            <Table.Cell>{item.msg}</Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table>
  )
}

export default LogList

/*
<Table>
  <Table.Header>
    <Table.Row>
      <Table.HeaderCell />
      <Table.HeaderCell />
    </Table.Row>
  </Table.Header>
  <Table.Body>
    <Table.Row>
      <Table.Cell />
      <Table.Cell />
    </Table.Row>
  </Table.Body>
  <Table.Footer />
</Table>
* */
