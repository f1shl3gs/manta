// Libraries
import React, {FunctionComponent} from 'react'
import {useNavigate} from 'react-router-dom'

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
import {useOrganizations} from './useOrganizations'

interface Props {
  visible: boolean
  dismiss: () => void
}

const OrganizationsSwitcher: FunctionComponent<Props> = ({
  visible,
  dismiss,
}) => {
  const {organizations, current, setCurrent} = useOrganizations()
  const navigate = useNavigate()

  return (
    <Overlay visible={visible}>
      <OverlayContainer maxWidth={500}>
        <OverlayHeader title="Switch Organization" onDismiss={dismiss} />

        <OverlayBody>
          <p className="org-switcher--prompt">Choose an organization</p>

          <List>
            {organizations.map(org => {
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
                    navigate(`/orgs/${org.id}/dashboards`)
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
