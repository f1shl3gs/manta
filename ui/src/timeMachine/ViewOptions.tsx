import React from 'react'

import {DapperScrollbars, Grid} from '@influxdata/clockface'

import OptionsSwitcher from 'src/timeMachine/OptionsSwitcher'
import {AppState} from 'src/types/stores'
import {useDispatch, useSelector} from 'react-redux'
import {setViewProperties} from 'src/timeMachine/actions'

const ViewOptions: React.FC = () => {
  const dispatch = useDispatch()

  const viewProperties = useSelector((state: AppState) => {
    return state.timeMachine.viewProperties
  })

  const update = viewProperties => {
    dispatch(setViewProperties(viewProperties))
  }

  return (
    <div className={'view-options'}>
      <DapperScrollbars
        autoHide={false}
        style={{width: '100%', height: '100%'}}
      >
        <div className={'view-options--container'}>
          <Grid>
            <Grid.Row>
              <OptionsSwitcher
                viewProperties={viewProperties}
                update={update}
              />
            </Grid.Row>
          </Grid>
        </div>
      </DapperScrollbars>
    </div>
  )
}

export default ViewOptions
