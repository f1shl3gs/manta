// Libraries
import React, {useState} from 'react'

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

const SigninForm: React.FC = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  return (
    <Form onSubmit={() => console.log('submit')}>
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
