import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'

import OptionsSwitcher from 'src/timeMachine/OptionsSwitcher'
import {AppState} from 'src/types/stores'
import {useSelector} from 'react-redux'

const ViewOptions: React.FC = () => {
  const viewProperties = useSelector((state: AppState) => {
    return state.timeMachine.viewProperties
  })

  return (
    <div className={'view-options'}>
      <DapperScrollbars
        autoHide={false}
        style={{width: '100%', height: '100%'}}
      >
        <div className={'view-options--container'}>
          <Grid>
            <Grid.Row>
              <OptionsSwitcher view={viewProperties} />
            </Grid.Row>
          </Grid>
        </div>
      </DapperScrollbars>
    </div>
  )
}

export default ViewOptions
