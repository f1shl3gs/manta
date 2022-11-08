// Libraries
import React, {FC, useState} from 'react'

// Components
import {
  Button,
  ButtonType,
  Columns,
  ComponentColor,
  ComponentSize,
  Form,
  FormElement,
  FormFooter,
  Grid,
  GridColumn,
  GridRow,
  Input,
  InputType,
} from '@influxdata/clockface'
import {useNavigate, useSearchParams} from 'react-router-dom'
import useFetch from 'src/shared/useFetch'

export const SignInForm: FC = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const returnTo = params.get('returnTo') || '/'

  const {run: submit} = useFetch(`/api/v1/signin`, {
    method: 'POST',
    body: {
      username,
      password,
    },
    onSuccess: () => navigate(returnTo),
    onError: err => {
      console.log('signin error', err)
    },
  })

  return (
    <Form onSubmit={() => submit()}>
      <Grid>
        <GridRow>
          <GridColumn widthXS={Columns.Twelve}>
            <FormElement label="Username">
              <Input
                name="username"
                value={username}
                onChange={ev => setUsername(ev.target.value)}
                size={ComponentSize.Medium}
                autoFocus={true}
                testID="username-input"
              />
            </FormElement>
          </GridColumn>

          <GridColumn widthXS={Columns.Twelve}>
            <FormElement label="Password">
              <Input
                name="password"
                testID={'password-input'}
                value={password}
                onChange={ev => setPassword(ev.target.value)}
                size={ComponentSize.Medium}
                type={InputType.Password}
              />
            </FormElement>
          </GridColumn>

          <GridColumn widthXS={Columns.Twelve}>
            <FormFooter>
              <Button
                color={ComponentColor.Primary}
                text="Sign In"
                size={ComponentSize.Medium}
                type={ButtonType.Submit}
                id="submit-signin"
                testID={'signin-button'}
              />
            </FormFooter>
          </GridColumn>
        </GridRow>
      </Grid>
    </Form>
  )
}
