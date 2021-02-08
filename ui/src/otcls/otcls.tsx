// Libraries
import React from 'react'
import moment from 'moment'

// Components
import {
  Button,
  ButtonShape,
  Columns,
  ComponentColor,
  ComponentSize,
  EmptyState,
  FlexBox,
  Grid,
  IconFont,
  ResourceCard,
  ResourceList,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import CopyButton from '../shared/components/CopyButton'

// Hooks
import {useOrgID} from 'shared/useOrg'
import {useOtcls} from './state'
import {useFetch} from 'use-http'
import {useHistory} from 'react-router-dom'

// Types
import {Otcl} from 'types/otcl'
import OtclExplainer from './component/OtclExplainer'

const Otcls: React.FC = () => {
  const orgID = useOrgID()
  const history = useHistory()
  const {otcls = [], rds, reload} = useOtcls()
  const {del} = useFetch(`/api/v1/otcls`, {})

  const context = (id: string): JSX.Element => {
    return (
      <FlexBox margin={ComponentSize.Small}>
        <FlexBox.Child>
          <Button
            icon={IconFont.Trash}
            text="Delete"
            size={ComponentSize.ExtraSmall}
            color={ComponentColor.Danger}
            onClick={() => {
              del(`${id}`)
                .then(() => {
                  console.log('delete', id, 'success')
                  reload()
                })
                .catch(() => {
                  console.log('failed')
                })
            }}
          />
        </FlexBox.Child>

        <FlexBox.Child>
          <CopyButton
            shape={ButtonShape.Default}
            textToCopy={`${window.location.protocol}//${window.location.hostname}/api/v1/otcls/${id}`}
            buttonText={'Copy Otcl config url'}
            color={ComponentColor.Default}
            contentName={'cn'}
            size={ComponentSize.ExtraSmall}
          />
        </FlexBox.Child>
      </FlexBox>
    )
  }

  return (
    <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
      <Grid>
        <Grid.Row>
          <Grid.Column
            widthXS={Columns.Twelve}
            widthSM={otcls.length !== 0 ? Columns.Eight : Columns.Twelve}
            widthMD={otcls.length !== 0 ? Columns.Ten : Columns.Twelve}
          >
            <ResourceList>
              <ResourceList.Body
                emptyState={
                  <EmptyState size={ComponentSize.Large}>
                    <EmptyState.Text>
                      Looks like this Org doesn't have any <b>Otcl</b> configs
                    </EmptyState.Text>
                  </EmptyState>
                }
              >
                {otcls?.map((item: Otcl) => (
                  <ResourceCard key={item.id} contextMenu={context(item.id)}>
                    <ResourceCard.Name
                      name={item.name}
                      onClick={() => {
                        history.push(`/orgs/${orgID}/otcls/${item.id}`)
                      }}
                    />
                    <ResourceCard.Description description={item.desc} />
                    <ResourceCard.Meta>
                      <span>updated: {moment(item.updated).fromNow()}</span>
                    </ResourceCard.Meta>
                  </ResourceCard>
                ))}
              </ResourceList.Body>
            </ResourceList>
          </Grid.Column>

          {otcls.length !== 0 && (
            <Grid.Column
              widthXS={Columns.Twelve}
              widthSM={Columns.Four}
              widthMD={Columns.Two}
            >
              <OtclExplainer />
            </Grid.Column>
          )}
        </Grid.Row>
      </Grid>
    </SpinnerContainer>
  )
}

export default Otcls
