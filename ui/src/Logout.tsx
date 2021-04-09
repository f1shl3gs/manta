import React, {useEffect} from 'react'
import {useHistory} from 'react-router-dom'

const Logout: React.FC = () => {
  const history = useHistory()

  useEffect(() => {
    fetch(`/api/v1/signout`, {
      method: 'DELETE',
    })
      .then(resp => {
        console.log('resp', resp)

        history.push(`/signin`)
      })
      .catch(err => {
        console.log('logout error: ' + err.message)
      })
  }, [history])

  return null
}

export default Logout
