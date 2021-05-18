// Libraries
import React, {useCallback, useState} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  Button,
  ButtonType,
  Columns,
  ComponentColor,
  ComponentSize,
  Form,
  Grid,
  Input,
  InputType,
} from '@influxdata/clockface'

// Hooks
import {useFetch} from 'shared/useFetch'
import {
  defaultErrorNotification,
  useNotification,
} from '../shared/notification/useNotification'

const SigninForm: React.FC = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const history = useHistory()
  const {notify} = useNotification()

  const {post, loading, error} = useFetch(`/api/v1/signin`, {})
  const onSubmit = useCallback(() => {
    post({
      username,
      password,
    })
      .then(() => {
        // success
        const params = new URLSearchParams(window.location.search)
        const returnTo = params.get('returnTo')
        if (!returnTo || returnTo === '') {
          history.push(`/orgs`)
        } else {
          history.push(`${decodeURIComponent(returnTo)}`)
        }
      })
      .catch(err => {
        notify({
          ...defaultErrorNotification,
          message: `Sign in failed, err: ${err.message}`,
        })
      })
  }, [post, username, password, history, notify])

  return (
    <Form onSubmit={onSubmit}>
      <Grid>
        <Grid.Row>
          <Grid.Column widthXS={Columns.Twelve}>
            <Form.Element label={'Username'}>
              <Input
                name={'username'}
                value={username}
                onChange={ev => setUsername(ev.target.value)}
                size={ComponentSize.Medium}
                autoFocus={true}
              />
            </Form.Element>
          </Grid.Column>

          <Grid.Column widthXS={Columns.Twelve}>
            <Form.Element label={'Password'}>
              <Input
                name={'password'}
                value={password}
                onChange={ev => setPassword(ev.target.value)}
                size={ComponentSize.Medium}
                type={InputType.Password}
              />
            </Form.Element>
          </Grid.Column>

          <Grid.Column widthXS={Columns.Twelve}>
            <Form.Footer>
              <Button
                color={ComponentColor.Primary}
                text={'Sign In'}
                size={ComponentSize.Medium}
                type={ButtonType.Submit}
              />
            </Form.Footer>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Form>
  )
}

export default SigninForm
