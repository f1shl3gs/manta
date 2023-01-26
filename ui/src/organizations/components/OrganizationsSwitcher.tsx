// Libraries
import React, {FunctionComponent, useCallback} from 'react'

// Components
import {
  ComponentSize,
  List,
  ListItem,
  Overlay,
  OverlayBody,
  OverlayContainer,
  OverlayHeader,
} from '@influxdata/clockface'

// Hooks
import {useNavigate} from 'react-router-dom'
import {useDispatch, useSelector} from 'react-redux'
import {getAll} from 'src/resources/selectors'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'
import {Organization} from 'src/types/organization'

// Actions
import {SET_ORG} from 'src/organizations/actions'

// Selectors
import {getOrg} from 'src/organizations/selectors'

interface Props {
  visible: boolean
  dismiss: () => void
}

const OrganizationsSwitcher: FunctionComponent<Props> = ({
  visible,
  dismiss,
}) => {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const {orgs, current} = useSelector((state: AppState) => {
    const orgs = getAll<Organization>(state, ResourceType.Organizations)
    const current = getOrg(state)

    return {
      orgs,
      current,
    }
  })
  const setCurrent = useCallback(
    org => {
      dispatch({
        type: SET_ORG,
        org,
      })
    },
    [dispatch]
  )

  return (
    <Overlay visible={visible}>
      <OverlayContainer maxWidth={500}>
        <OverlayHeader title="Switch Organization" onDismiss={dismiss} />

        <OverlayBody>
          <p className="org-switcher--prompt">Choose an organization</p>

          <List>
            {orgs.map(org => {
              const selected = org === current

              return (
                <ListItem
                  key={org.id}
                  size={ComponentSize.Large}
                  selected={selected}
                  wrapText={false}
                  onClick={() => {
                    if (selected) {
                      dismiss()
                      return
                    }

                    setCurrent(org)
                    dismiss()
                    navigate(`/orgs/${org.id}`)
                  }}
                >
                  {org.name}
                </ListItem>
              )
            })}
          </List>
        </OverlayBody>
      </OverlayContainer>
    </Overlay>
  )
}

export default OrganizationsSwitcher
