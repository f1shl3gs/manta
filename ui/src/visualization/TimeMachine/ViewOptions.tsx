import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'

import OptionsSwitcher from 'src/visualization/TimeMachine/OptionsSwitcher'
import {useViewProperties} from 'src/visualization/TimeMachine/useViewProperties'

const ViewOptions: React.FC = () => {
  const {viewProperties} = useViewProperties()

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
