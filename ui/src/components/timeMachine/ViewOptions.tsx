import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'

import {ViewProperties} from 'types'
import OptionsSwitcher from './OptionsSwitcher'

interface Props {
  view: ViewProperties
}

const ViewOptions: React.FC<Props> = (props) => {
  const {view} = props

  return (
    <div className={'view-options'}>
      <DapperScrollbars
        autoHide={false}
        style={{width: '100%', height: '100%'}}
      >
        <div className={'view-options--container'}>
          <Grid>
            <Grid.Row>
              <OptionsSwitcher view={view} />
            </Grid.Row>
          </Grid>
        </div>
      </DapperScrollbars>
    </div>
  )
}

export default ViewOptions
