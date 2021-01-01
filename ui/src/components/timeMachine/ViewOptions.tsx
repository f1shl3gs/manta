import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'

import OptionsSwitcher from './OptionsSwitcher'
import { useViewProperties } from './useView';

interface Props {

}

const ViewOptions: React.FC<Props> = (props) => {
  const {viewProperties} = useViewProperties();

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
