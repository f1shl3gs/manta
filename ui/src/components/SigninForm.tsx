// Libraries
import React, {useCallback, useState} from 'react'

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
import {useFetch} from 'use-http'
import {useHistory, useLocation} from 'react-router-dom'

const SigninForm: React.FC = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const history = useHistory()
  const location = useLocation()

  const {post, loading, error} = useFetch(`/api/v1/signin`, {})
  const onSubmit = useCallback(() => {
    post({
      username,
      password,
    })
      .then(() => {
        // success
        const q = new URLSearchParams(location.search)
        history.push(`${decodeURIComponent(q.get('returnTo') || '/')}`)
      })
      .catch(() => {
        // failed
      })
  }, [username, password])

  return (
    <Form onSubmit={onSubmit}>
      <Grid>
        <Grid.Row>
          <Grid.Column widthXS={Columns.Twelve}>
            <Form.Element label={'Username'}>
              <Input
                name={'username'}
                value={username}
                onChange={(ev) => setUsername(ev.target.value)}
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
                onChange={(ev) => setPassword(ev.target.value)}
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
