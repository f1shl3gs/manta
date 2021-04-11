import React, {useEffect} from 'react'
import {useHistory} from 'react-router-dom'

const Logout: React.FC = () => {
  const history = useHistory()

  useEffect(() => {
    fetch(`/api/v1/signout`, {
      method: 'DELETE',
    })
      .then(resp => {
        if (resp.status === 204) {
          history.push(`/signin`)
          return
        }

        // todo: handle failures
        console.log('unexpected response', resp)
      })
      .catch(err => {
        console.log('logout error: ' + err.message)
      })
  }, [history])

  return null
}

export default Logout
