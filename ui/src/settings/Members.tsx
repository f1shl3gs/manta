import React, {FunctionComponent} from 'react'
import {GetResources, ResourceType} from 'src/shared/components/GetResources'
import MemberList from 'src/settings/MemberList'

const Members: FunctionComponent = () => {
  return (
    <GetResources type={ResourceType.Users}>
      <MemberList />
    </GetResources>
  )
}

export default Members
