// Libraries
import React from 'react'

// Components
import {ComponentColor, ComponentSize, Form, Grid} from '@influxdata/clockface'
import DashedButton from 'shared/components/DashedButton'
import DurationInput from 'shared/components/DurationInput'
import BuilderCard from '../builderCard/BuilderCard'
import BuilderCardHeader from '../builderCard/BuilderCardHeader'
import BuilderCardBody from '../builderCard/BuilderCardBody'
import CheckLabelRow from './CheckLabelRow'

// Hooks
import {useCheck} from '../checks/useCheck'

// Constants
import {DURATIONS} from 'constants/duration'

const CheckMetaCard: React.FC = () => {
  const {labels} = useCheck()

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

        <Form.Label label={'Labels'} />

        {labels.map((label, index) => (
          <CheckLabelRow
            key={index}
            index={index}
            label={label}
            handleChangeLabelRow={() => console.log('change')}
            handleRemoveTagRow={() => console.log('remove')}
          />
        ))}

        <DashedButton
          text={'+ Label'}
          onClick={() => console.log('add label')}
          color={ComponentColor.Primary}
          size={ComponentSize.Small}
        />
      </BuilderCardBody>
    </BuilderCard>
  )
}

export default CheckMetaCard
