// Libraries
import React, {ChangeEvent, FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {Columns, Form, Grid, Input} from '@influxdata/clockface'
import EndpointTypeDropdown from 'src/notification_endpoints/components/EndpointTypeDropdown'
import EndpointOptions from 'src/notification_endpoints/components/EndpointOptions'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {updateNotificationEndpoint} from 'src/notification_endpoints/actions/creators'

// Selectors
import {getEndpoint} from 'src/notification_endpoints/selectors'

const mdtp = {
  onUpdate: updateNotificationEndpoint,
}

const mstp = (state: AppState) => {
  return {
    endpoint: getEndpoint(state),
  }
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const NotificationEndpointForm: FunctionComponent<Props> = ({
  endpoint,
  onUpdate,
}) => {
  const handleChange = (ev: ChangeEvent<HTMLInputElement>) => {
    const {name, value} = ev.target

    onUpdate({
      ...endpoint,
      [name]: value,
    })
  }

  return (
    <Form>
      <Grid>
        <Grid.Row>
          <Grid.Column widthSM={Columns.Six}>
            <Form.Element label={'Destination'}>
              <EndpointTypeDropdown
                selected={endpoint.type}
                onSelect={type => console.log('select type', type)}
              />
            </Form.Element>
          </Grid.Column>

          <Grid.Column widthSM={Columns.Six}>
            <Form.Element label={'Name'}>
              <Input
                name={'name'}
                testID={'endpoint-name--input'}
                placeholder={'Name this endpoint'}
                value={endpoint.name}
                onChange={handleChange}
              />
            </Form.Element>
          </Grid.Column>
        </Grid.Row>

        <Grid.Row>
          <Grid.Column>
            <Form.Element label={'Description'}>
              <Input
                name={'desc'}
                testID={'endpoint-desc--input'}
                placeholder={'Describe this endpoint'}
                value={endpoint.desc}
                onChange={handleChange}
              />
            </Form.Element>
          </Grid.Column>
        </Grid.Row>

        <Grid.Row>
          <Grid.Column widthSM={Columns.Twelve}>
            <EndpointOptions type={endpoint.type} />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Form>
  )
}

export default connector(NotificationEndpointForm)
