// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import BuilderCard from 'src/checks/components/builderCard/BuilderCard'
import {
  ComponentColor,
  FormElement,
  Grid,
  GridColumn,
  GridRow,
  Input,
  QuestionMarkTooltip,
} from '@influxdata/clockface'
import BuilderCardHeader from 'src/checks/components/builderCard/BuilderCardHeader'
import BuilderCardBody from 'src/checks/components/builderCard/BuilderCardBody'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {setCron} from 'src/checks/actions/builder'

const mstp = (state: AppState) => {
  const {cron, conditions} = state.checkBuilder

  return {
    cron,
    conditions,
  }
}

const mdtp = {
  setCron,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const PropertiesCard: FunctionComponent<Props> = ({cron, setCron}) => {
  const tooltipContents = (
    <>
      <strong>Schedule</strong> defines when it should be scheduled,
      <br/>
      Manta use <strong>cron</strong> to implement this feature, so
      value like " */5 * * * * " is supported.
      <br/>
      <br/>
      For convienent, <strong>@every 1m</strong> is supported too.
    </>
  )

  return (
    <BuilderCard
      testID={'builder-card-properties'}
      className="alert-builder--card alert-builder--meta-card"
    >
      <BuilderCardHeader title={'Properties'} />
      <BuilderCardBody autoHideScrollbars={true} addPadding={true}>
        <Grid>
          <GridRow>
            <GridColumn>
              <FormElement
                label={'Schedule'}
                labelAddOn={() => (
                  <QuestionMarkTooltip
                    diameter={18}
                    color={ComponentColor.Primary}
                    testID={`aaa--question-mark`}
                    tooltipContents={tooltipContents}
                  />
                )}
              >
                <Input value={cron} onChange={ev => setCron(ev.target.value)} />
              </FormElement>
            </GridColumn>
          </GridRow>
        </Grid>
      </BuilderCardBody>
    </BuilderCard>
  )
}

export default connector(PropertiesCard)
