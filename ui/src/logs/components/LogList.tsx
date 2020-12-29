// libraries
import React, { useMemo } from 'react';
import moment from 'moment';

// components
import { Config, DEFAULT_TABLE_COLORS, fromRows, HoverTimeProvider, Plot } from '@influxdata/giraffe';
import { TableGraphLayerConfig } from '@influxdata/giraffe/dist/types';

// test data
import testData from './query_range_resp.json';

const tableCSV = `#group,false,false,true,true,false,false,true,true,true,true
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#default,_result,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,cpu,host
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:45:50Z,2.2,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:00Z,1,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:10Z,1.2,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:20Z,0.8,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:30Z,0.8008008008008008,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:40Z,0.6993006993006993,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:50Z,1.001001001001001,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:00Z,0.6993006993006993,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:10Z,0.7007007007007007,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:20Z,0.8,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:30Z,0.7992007992007992,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:40Z,0.8,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:50Z,0.8,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:00Z,0.9,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:10Z,0.9,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:20Z,0.8008008008008008,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:30Z,0.8991008991008991,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:40Z,0.7,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:50Z,0.9,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:00Z,0.7007007007007007,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:10Z,0.7992007992007992,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:20Z,0.9009009009009009,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:30Z,1.098901098901099,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:40Z,0.8,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:50Z,0.9,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:00Z,0.9,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:10Z,0.7,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:20Z,0.7,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,0,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:30Z,0.9009009009009009,usage_system,cpu,cpu1,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:45:50Z,21.678321678321677,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:00Z,12.2,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:10Z,19.51951951951952,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:20Z,11.9,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:30Z,13.386613386613387,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:40Z,11.4,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:46:50Z,12.1,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:00Z,10.91091091091091,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:10Z,10.589410589410589,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:20Z,12.612612612612613,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:30Z,10.289710289710289,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:40Z,11.7,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:47:50Z,10.5,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:00Z,10.6,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:10Z,13.3,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:20Z,13.6,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:30Z,11.511511511511511,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:40Z,11.3,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:48:50Z,12.8,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:00Z,10.989010989010989,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:10Z,10.71071071071071,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:20Z,11.711711711711711,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:30Z,13.086913086913087,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:40Z,10.989010989010989,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:49:50Z,11.511511511511511,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:00Z,11.488511488511488,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:10Z,10.6,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:20Z,10.21021021021021,usage_system,cpu,cpu0,MBP15-TLUONG.local
,,1,2020-07-01T21:45:43.0968Z,2020-07-01T21:50:43.0968Z,2020-07-01T21:50:30Z,11.4,usage_system,cpu,cpu0,MBP15-TLUONG.local
`;
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
  const rows = new Array<TableRow>();

  ss.forEach((stream) => {
    stream.values.forEach((pair) => {
      rows.push({
        ts: moment(Number(pair[0]) / 1000 / 1000).format(),
        msg: pair[1]
      });
    });
  });

  return rows;
};

interface Result {
  stream: {
    [key: string]: string
  }
  values: [
    [string, string]
  ]
}

// Notice: it is worked now, i don't know what happened
// and TableGraphLayerConfig is not export at @influxdata/girrafe (through it can be access by it's dist/types)
// maybe it will works in the future
const LogList = () => {
  // todo: show the common labels
  /*

    const resp = testData as Resp;

    // transformData
    const result = resp.data.result as StreamResult[];
    const rows = result.map((result: StreamResult) => {
      const { stream, values } = result;

      return values.map(val => {
        return {
          ...stream,
          time: val[0],
          log: val[1]
        };
      });
    }).flat();

    const table = fromRows(rows);
    console.log('table', table);
  */

  const theme = 'dark';
  const fixFirstColumn = false;
  const config: Config = {
    fluxResponse: tableCSV,
    layers: [
      {
        type: 'table',
        properties: {
          colors: DEFAULT_TABLE_COLORS,
          tableOptions: {
            fixFirstColumn,
            verticalTimeAxis: true
          },
          fieldOptions: [
            {
              displayName: '_start',
              internalName: '_start',
              visible: true
            },
            {
              displayName: '_stop',
              internalName: '_stop',
              visible: true
            },
            {
              displayName: '_time',
              internalName: '_time',
              visible: true
            },
            {
              displayName: '_value',
              internalName: '_value',
              visible: true
            },
            {
              displayName: '_field',
              internalName: '_field',
              visible: true
            },
            {
              displayName: '_measurement',
              internalName: '_measurement',
              visible: true
            },
            {
              displayName: 'cpu',
              internalName: 'cpu',
              visible: true
            },
            {
              displayName: 'host',
              internalName: 'host',
              visible: true
            }
          ],
          timeFormat: 'YYYY/MM/DD HH:mm:ss',
          decimalPlaces: {
            digits: 3,
            isEnforced: true
          }
        },
        timeZone: 'Local',
        tableTheme: theme
      } as TableGraphLayerConfig
    ]
  };

  return (
    <HoverTimeProvider>
      <div
        style={{
          width: 'calc(100vw - 100px)',
          height: 'calc(100vh - 100px)',
          margin: '50px'
        }}
      >
        <Plot config={config} />
      </div>
    </HoverTimeProvider>
  );
};

export default LogList;
