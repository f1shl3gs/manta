// Libraries
import React from 'react'

import {BorderType, ComponentSize, Table} from '@influxdata/clockface'
import testData from './testData'
import {colorBasedOnPackageName} from './color'
import Color from 'color'

interface Props {}

const generateTable = (): {name: string; self: number; total: number}[] => {
  const flamebearer = testData.flamebearer

  const {names, levels} = flamebearer
  const hash = {} as {
    [key: string]: {name: string; self: number; total: number}
  }

  for (let i = 0; i < levels.length; i++) {
    for (let j = 0; j < levels[i].length; j += 4) {
      const key = levels[i][j + 3]
      const name = names[key]
      hash[name] = hash[name] || {
        name: name || '<empty>',
        self: 0,
        total: 0,
      }

      hash[name].total += levels[i][j + 1]
      hash[name].self += levels[i][j + 2]
    }
  }

  return Object.values(hash)
}

const second = 1
const minute = 60 * second
const hour = 60 * minute
const day = 24 * hour

const durationUnits = ['seconds', 'minutes', 'hours', 'days', 'weeks']

const durationFormatter = (s: number): string => {
  s /= 100

  if (s < 0.01) {
    return '< 0.01'
  }

  if (s < minute) {
    return `${s.toFixed(2)} seconds`
  }

  s /= minute
  if (s < hour) {
    return `${s.toFixed(2)} minutes`
  }

  s /= hour
  if (s < day) {
    return `${s.toFixed(2)} hours`
  }

  return `${s.toFixed(2)} weeks`
}

export function getPackageNameFromStackTrace(
  spyName: string,
  stackTrace: string
) {
  // TODO: actually make sure these make sense and add tests
  const regexpLookup = {
    pyspy: /^(?<packageName>(.*\/)*)(?<filename>.*\.py+)(?<line_info>.*)$/,
    rbspy: /^(?<packageName>(.*\/)*)(?<filename>.*\.rb+)(?<line_info>.*)$/,
    gospy: /^(?<packageName>(.*\/)*)(?<filename>.*)(?<line_info>.*)$/,
    ebpfspy: /^(?<packageName>.+)$/,
    default: /^(?<packageName>(.*\/)*)(?<filename>.*)(?<line_info>.*)$/,
  } as {
    [key: string]: RegExp
  }

  if (stackTrace.length === 0) {
    return stackTrace
  }
  const regexp = regexpLookup[spyName] || regexpLookup.default
  const fullStackGroups = stackTrace.match(regexp)
  if (fullStackGroups) {
    // @ts-ignore
    return fullStackGroups.groups.packageName
  }
  return stackTrace
}

const backgroundImageStyle = (a: number, b: number, color: Color) => {
  const w = 148
  const k = w - (a / b) * w
  const clr = color.alpha(1.0)

  return {
    backgroundImage: `linear-gradient(${clr}, ${clr})`,
    backgroundPosition: `-${k}px 0px`,
    backgroundRepeat: 'no-repeat',
    width: '148px',
    color: 'white',
  }
}

const TablePanel: React.FC<Props> = props => {
  const table = generateTable()
  const {spyName, maxSelf, numTicks} = testData.flamebearer

  const sorted = table.sort((a, b) => b.self - a.self)

  return (
    <Table
      key={'table-panel'}
      fontSize={ComponentSize.ExtraSmall}
      striped={true}
      borders={BorderType.All}
    >
      <Table.Header>
        <Table.Row>
          <Table.HeaderCell>Location</Table.HeaderCell>
          <Table.HeaderCell>Self</Table.HeaderCell>
          <Table.HeaderCell>Total</Table.HeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {sorted.map(item => {
          const pn = getPackageNameFromStackTrace(spyName, item.name)
          const color = colorBasedOnPackageName(pn, 0.9)
          const style = {
            backgroundColor: color,
          }

          return (
            <Table.Row key={item.name}>
              <Table.Cell style={{color: 'white'}}>
                {/*  @ts-ignore */}
                <span className="color-reference" style={style} />
                <span>{item.name}</span>
              </Table.Cell>
              <Table.Cell
                className={'profile-table-cell'}
                style={backgroundImageStyle(item.self, maxSelf, color)}
              >
                {durationFormatter(item.self)}
              </Table.Cell>
              <Table.Cell
                className={'profile-table-cell'}
                style={backgroundImageStyle(item.total, numTicks, color)}
              >
                {durationFormatter(item.total)}
              </Table.Cell>
            </Table.Row>
          )
        })}
      </Table.Body>
    </Table>
  )
}

export default TablePanel
