import React, {useState} from 'react'
import {
  Button,
  ComponentSize,
  DapperScrollbars,
  FlexBox,
  FlexDirection,
  IconFont,
  InfluxColors,
  Input,
  Panel,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import CheckCards from './CheckCards'
import {useChecks} from './useChecks'

const ChecksColumn: React.FC = () => {
  const {checks, remoteDataState} = useChecks()
  const title = `Checks`
  const [search, setSearch] = useState('')

  const tooltipContents = (
    <>
      A <strong>Check</strong> is a periodic query that the system
      <br />
      performs against your time series data
      <br />
      that will generate a status
      <br />
      <br />
    </>
  )

  const createButton = <Button titleText={'Create'} text={'Create'} />
  const panelClassName = `alerting-index--column alerting-index--checks`

  return (
    <Panel backgroundColor={InfluxColors.Kevlar} className={panelClassName}>
      <Panel.Header>
        <FlexBox direction={FlexDirection.Row} margin={ComponentSize.Small}>
          <h4 style={{width: 'auto', marginRight: '6px'}}>{title}</h4>
          {createButton}
        </FlexBox>
      </Panel.Header>
      <div className={'alerting-index--search'}>
        <Input
          icon={IconFont.Search}
          placeholder={`Filter ${title}`}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>
      <div className={'alerting-index--column-body'}>
        <DapperScrollbars
          autoHide={true}
          style={{width: '100%', height: '100%'}}
        >
          <div className={'alerting-index--list'}>
            <SpinnerContainer
              loading={remoteDataState}
              spinnerComponent={<TechnoSpinner />}
            >
              <CheckCards search={''} checks={checks} />
            </SpinnerContainer>
          </div>
        </DapperScrollbars>
      </div>
    </Panel>
  )
}

export default ChecksColumn
