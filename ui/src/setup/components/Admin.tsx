// Library
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  Button,
  Columns,
  ComponentColor,
  ComponentStatus,
  Form,
  Grid,
  Input,
  InputType,
} from '@influxdata/clockface'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {initialUser} from 'src/setup/actions/thunks'
import {setSetupParams} from 'src/setup/actions/creators'

const mstp = (state: AppState) => {
  const {username, password, organization} = state.setup

  return {
    username,
    password,
    organization,
  }
}

const mdtp = {
  submit: initialUser,
  updateParams: setSetupParams,
}

const connector = connect(mstp, mdtp)

type Props = ConnectedProps<typeof connector>

const Admin: FunctionComponent<Props> = ({
  username,
  password,
  organization,
  submit,
  updateParams,
}) => {
  const submitStatus =
    username !== '' && password !== ''
      ? ComponentStatus.Default
      : ComponentStatus.Disabled

  const setUsername = ev => {
    updateParams({
      username: ev.target.value,
    })
  }

  const setPassword = ev =>
    updateParams({
      password: ev.target.value,
    })

  const setOrganization = ev =>
    updateParams({
      organization: ev.target.value,
    })

  return (
    <div className={'onboarding-step'}>
      <Form>
        <div className="wizard-step--scroll-area">
          <div className="wizard-step--scroll-content">
            <h3 className="wizard-step--title">Setup Initial User</h3>

            <h5 className={'wizard-step--sub-title'}>
              You will be able to create additional Users and Organizations
              later
            </h5>

            <Grid>
              <Grid.Row>
                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthMD={Columns.Ten}
                  offsetMD={Columns.One}
                >
                  <Form.Element label="Username">
                    <Input
                      testID={'input-username'}
                      autoFocus={true}
                      value={username}
                      onChange={setUsername}
                    />
                  </Form.Element>
                </Grid.Column>

                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthMD={Columns.Ten}
                  offsetMD={Columns.One}
                >
                  <Form.Element label="Password">
                    <Input
                      testID={'input-password'}
                      type={InputType.Password}
                      value={password}
                      onChange={setPassword}
                    />
                  </Form.Element>
                </Grid.Column>

                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthMD={Columns.Ten}
                  offsetMD={Columns.One}
                >
                  <Form.Element label="Organization">
                    <Input
                      testID={'input-organization'}
                      value={organization}
                      onChange={setOrganization}
                    ></Input>
                  </Form.Element>
                </Grid.Column>

                <Grid.Column
                  widthXS={Columns.Twelve}
                  widthMD={Columns.Ten}
                  offsetMD={Columns.Five}
                >
                  <Button
                    testID={'button-next'}
                    text="Next"
                    color={ComponentColor.Primary}
                    status={submitStatus}
                    onClick={() => {
                      submit()
                    }}
                  />
                </Grid.Column>
              </Grid.Row>
            </Grid>
          </div>
        </div>
      </Form>
    </div>
  )
}

export default connector(Admin)
