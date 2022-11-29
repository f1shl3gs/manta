import React, {FunctionComponent, useEffect} from 'react'
import {useNavigate} from 'react-router-dom'
import {useOrg} from './selectors'

const ToOrg: FunctionComponent = () => {
  const navigate = useNavigate()
  const {id} = useOrg()

  useEffect(() => {
    console.log('toorg')
    navigate(`/orgs/${id}`)
  }, [id, navigate])

  return <></>
}

export default ToOrg
