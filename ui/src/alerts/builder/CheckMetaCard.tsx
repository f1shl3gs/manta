// Libraries
import React from 'react'

// Components
import BuilderCard from '../builderCard/BuilderCard'
import BuilderCardHeader from '../builderCard/BuilderCardHeader'
import BuilderCardBody from '../builderCard/BuilderCardBody'
import {Form, Grid} from '@influxdata/clockface'
import DurationInput from '../../shared/components/DurationInput'

// Constants
import {DURATIONS} from '../../constants/duration'

const CheckMetaCard: React.FC = () => {
  return (
    <BuilderCard
      testID={'builder-meta'}
      className={'alert-builder--card alert-builder--meta-card'}
    >
      <BuilderCardHeader title={'Properties'} />
      <BuilderCardBody addPadding={true} autoHideScrollbars={true}>
        <Grid>
          <Grid.Row>
            <Grid.Column widthSM={6}>
              <Form.Element label={'Schedule Every'}>
                <DurationInput
                  value={'5m'}
                  suggestions={DURATIONS}
                  onSubmit={v => console.log(v)}
                  testID={'schedule-check'}
                />
              </Form.Element>
            </Grid.Column>

            <Grid.Column widthSM={6}>
              <Form.Element label="Offset">
                <DurationInput
                  value={'1m'}
                  suggestions={DURATIONS}
                  onSubmit={v => console.log(v)}
                  testID="offset-options"
                />
              </Form.Element>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </BuilderCardBody>
    </BuilderCard>
  )
}

export default CheckMetaCard
