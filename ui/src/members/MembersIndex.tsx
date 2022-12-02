import React, {FunctionComponent} from 'react'
import MemberList from 'src/members/MemberList'
import GetResources from 'src/resources/components/GetResources'
import {ResourceType} from 'src/types/resources'

const MembersIndex: FunctionComponent = () => {
  return (
    <GetResources resources={[ResourceType.Users]}>
      <MemberList />
    </GetResources>
  )
}

export default MembersIndex
