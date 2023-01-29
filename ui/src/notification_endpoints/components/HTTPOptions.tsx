// Libraries
import React, {ChangeEvent, FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {Columns, Form, Grid, Input, Panel} from '@influxdata/clockface'
import MethodDropdown from 'src/notification_endpoints/components/MethodDropdown'
import AuthMethodDropdown from 'src/notification_endpoints/components/AuthMethodDropdown'

// Selectors
import {getEndpoint} from 'src/notification_endpoints/selectors'
import {AppState} from '../../types/stores'
import {updateNotificationEndpoint} from '../actions/creators'

const mstp = (state: AppState) => {
  return {
    endpoint: getEndpoint(state),
  }
}

const mdtp = {
  onUpdate: updateNotificationEndpoint,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const HTTPOptions: FunctionComponent<Props> = ({endpoint, onUpdate}) => {
  const handleChange = (ev: ChangeEvent<HTMLInputElement>) => {
    const {name, value} = ev.target

    onUpdate({
      ...endpoint,
      [name]: value,
    })
  }

  return (
    <Panel>
      <Panel.Header>
        <h4>HTTP Options</h4>
      </Panel.Header>

      <Panel.Body>
        <Grid>
          <Grid.Row>
            <Grid.Column widthSM={Columns.Six}>
              <Form.Element label={'HTTP Method'}>
                <MethodDropdown
                  selected={'POST'}
                  onSelect={m => console.log('method', m)}
                />
              </Form.Element>
            </Grid.Column>
            <Grid.Column widthSM={Columns.Six}>
              <Form.Element label={'Auth Method'}>
                <AuthMethodDropdown
                  selected={endpoint.authMethod}
                  onSelect={m => console.log(m)}
                />
              </Form.Element>
            </Grid.Column>
          </Grid.Row>

          <Grid.Row>
            <Grid.Column>
              <Form.Element label={'URL'}>
                <Input
                  name={'url'}
                  value={endpoint.url}
                  placeholder={'http or https'}
                  testID={'http-options-url--input'}
                  onChange={handleChange}
                />
              </Form.Element>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Panel.Body>
    </Panel>
  )
}

export default connector(HTTPOptions)
