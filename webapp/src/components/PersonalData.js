import { useState, useEffect, useRef } from "react";
import { useTranslation } from 'react-i18next';
import UserService from '../services/user';
import GroupService from '../services/group';

const PersonalData = () => {
  const { t } = useTranslation(['translation', 'profile']);
  const groupTitleInput = useRef(null);
  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [email, setEmail] = useState("")
  const [isTeacher, setIsTeacher] = useState(false)
  const [ownedGroups, setOwnedGroups] = useState([])
  const [sharedGroups, setSharedGroups] = useState([])

  useEffect(()=> {
    UserService.getProfile()
    .then(
      (response) => {
        setFirstName(response.data.FirstName)
        setLastName(response.data.LastName)
        setEmail(response.data.Email)
        setIsTeacher(response.data.IsTeacher)
        if (response.data.OwnedGroups.Groups) {
          setOwnedGroups(response.data.OwnedGroups.Groups)
        }

        if (response.data.SharedGroups.Groups) {
          setSharedGroups(response.data.SharedGroups.Groups)
        }
      },
      (error) => {
        console.log(error.message)
      }
    )
  }, [])

  const handleCreateGroup = () => {
    GroupService.createGroup(groupTitleInput.current.value)
    .then(
      (response) => {
        window.location.reload();
      },
      (error) => {
        console.log(error.message)
      }
    )
  }

  return(
      <div>
        <div>
            <p className='text-purple-900'>
            <strong>{t('profile:firstname', 'First Name:')} </strong> {firstName}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:lastname', 'Last Name:')}  </strong> {lastName}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:email', 'Email:')} </strong> {email}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:role', 'Role:')} </strong> {isTeacher ? t('profile:teacher', 'Teacher') : t('profile:student', 'Student')}
          </p>
        </div>

        <div className="flex items-center justify-between mt-2">
            <input
              ref={groupTitleInput}
              className='relative block w-2/3 px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
            />
            <button
              className='relative flex justify-center px-2 py-1 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
              onClick={handleCreateGroup}
            >
              {t('profile:createGroup', 'Create Group')}
            </button>
        </div>

        { isTeacher ? (
          <div>
            <p className='mt-3 text-purple-900'>
              <strong>Owned Groups</strong>
            </p>
            <ul>
              {ownedGroups.map((group, idx) => {
                return <li className="hover:text-gray-600"><a href={`/groups/${group.ID}`}>{group.Title}</a></li>
              })}

              {ownedGroups === 0 ? <p>{t('profile:nogroups', "You have no owned groups")}</p> : null }
            </ul>
          </div>
        ) : 
          null
        }

        <div>
          <p className='mt-3 text-purple-900'>
            <strong>Shared Groups</strong>
          </p>
          <ul>
            {sharedGroups.map((group, idx) => {
              return <li className="hover:text-gray-600"><a href={`/groups/${group.ID}`}>{group.Title}</a></li>
            })}

            {sharedGroups.length === 0 ? <p>{t('profile:nogroups', "You have no shared groups")}</p> : null }
          </ul>
        </div>
      </div>
  )
}

export default PersonalData;