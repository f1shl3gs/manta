// Library
import React, {FC} from 'react'

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

// Hooks
import {useOnboard} from '../useOnboard'

export const Admin: FC = () => {
  const {
    username,
    password,
    setUsername,
    setPassword,
    organization,
    setOrganization,
    onboard,
  } = useOnboard()
  const submitStatus =
    username !== '' && password !== ''
      ? ComponentStatus.Default
      : ComponentStatus.Disabled

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
                      onChange={e => setUsername(e.target.value)}
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
                      onChange={e => setPassword(e.target.value)}
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
                      onChange={e => setOrganization(e.target.value)}
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
                    onClick={onboard}
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
