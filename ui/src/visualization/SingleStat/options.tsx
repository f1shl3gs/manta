//
import React, {FunctionComponent} from 'react'
import {SingleStatViewProperties} from 'src/types/cells'
import {Color, Columns, Form, Grid, Input} from '@influxdata/clockface'
import DecimalPlaces from '../components/DecimalPlaces'

interface Props {
  viewProperties: SingleStatViewProperties
  update: (viewProperties: any) => void
}

const SingleStatOptions: FunctionComponent<Props> = ({
  viewProperties,
  update,
}) => {
  const _setColors = (colors: Color[]): void => {
    update({colors})
  }

  return (
    <Grid>
      <Grid.Row>
        <Grid.Column>
          <Grid.Row>
            <Grid.Column widthSM={Columns.Six}>
              <Form.Element label="Prefix">
                <Input
                  value={viewProperties.prefix}
                  placeholder=""
                  onChange={ev => update({prefix: ev.target.value})}
                />
              </Form.Element>
            </Grid.Column>
            <Grid.Column widthSM={Columns.Six}>
              <Form.Element label="Suffix">
                <Input
                  value={viewProperties.suffix}
                  placeholder="%, rpm, etc."
                  onChange={ev => update({suffix: ev.target.value})}
                />
              </Form.Element>
            </Grid.Column>
          </Grid.Row>
        </Grid.Column>

        <Grid.Column>
          <DecimalPlaces
            isEnforced={viewProperties?.decimalPlaces?.isEnforced === true}
            digits={
              typeof viewProperties?.decimalPlaces?.digits === 'number' ||
              viewProperties?.decimalPlaces?.digits === null
                ? viewProperties.decimalPlaces.digits
                : NaN
            }
            update={update}
          />
        </Grid.Column>
      </Grid.Row>
    </Grid>
  )
}

export default SingleStatOptions
