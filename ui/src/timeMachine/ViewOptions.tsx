import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'
import OptionsSwitcher from 'src/timeMachine/OptionsSwitcher'

const ViewOptions: React.FC = () => {

  return (
    <div className={'view-options'}>
      <DapperScrollbars
        autoHide={false}
        style={{width: '100%', height: '100%'}}
      >
        <div className={'view-options--container'}>
          <Grid>
            <Grid.Row>
              <OptionsSwitcher />
            </Grid.Row>
          </Grid>
        </div>
      </DapperScrollbars>
    </div>
  )
}

export default ViewOptions
