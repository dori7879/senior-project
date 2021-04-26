import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import MultiSelect from "react-multi-select-component";

import Footer from '../components/Footer'
import Header from '../components/Header'
import { Redirect, useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import GroupService from '../services/group'
import UserService from '../services/user'
import { BASE_API_URL } from '../services'

const Group = () => {
  const { groupID } = useParams()
  const { t } = useTranslation(['translation', 'profile']);
  // eslint-disable-next-line no-unused-vars
  const { token, role } = useSelector(
    (state) => state.auth
  )
  // eslint-disable-next-line no-unused-vars
  const [title, setTitle] = useState("")
  const [shareLink, setShareLink] = useState("")
  const [ownerFirstName, setOwnerFirstName] = useState("")
  const [ownerLastName, setOwnerLastName] = useState("")
  const [ownerEmail, setOwnerEmail] = useState("")
  const [teachers, setTeachers] = useState([])
  const [members, setMembers] = useState([])
  const [selected, setSelected] = useState([])
  const [userOptions, setUserOptions] = useState([])


  const filterOptions = async (options, filter) => {
    if (!filter) {
      return options
    }
    
    let response = await UserService.fetchUserSuggestions(filter, false)
     
    let reqUsers = response.data.Users.map((u, i) => {
      return {
        // label: u.FirstName + " <" + u.Email + "> " + u.LastName,
        label: u.FirstName + " " + u.LastName,
        value: u.ID.toString(), 
      }
    })

    // Persist selected options
    let selectedOptions = options.filter((obj) => selected.includes(obj))

    // Remove users already in the group
    let memberIDs = members.map((u, i) => u.ID)
    let newOptions = reqUsers.filter((obj) => {
      return !memberIDs.includes(+obj.value)
    })
    
    return [...selectedOptions, ...newOptions];
  };

  const onSubmit = () => {
    let memberIDs = selected.map((obj, i) => +obj.value)
    GroupService.addMembers(groupID, memberIDs)
      .then(
        (response) => {
          setSelected([])
          setUserOptions([])
          window.location.reload()
        },
        (error) => {
          console.log(error.message)
        }
      )
  };

  useEffect(()=> {
    GroupService.getGroup(groupID)
    .then(
      (response) => {
        setTitle(response.data.Title)
        setShareLink(response.data.ShareLink)
        setOwnerEmail(response.data.Owner.Email)
        setOwnerFirstName(response.data.Owner.FirstName)
        setOwnerLastName(response.data.Owner.LastName)
        if (response.data.Teachers.Users) {
          setTeachers(response.data.Teachers.Users)
        }
        if (response.data.Members.Users) {
          setMembers(response.data.Members.Users)
        }
      },
      (error) => {
        console.log(error.message)
      }
    )
  }, [groupID])

  if (!token) {
    return <Redirect to='/login' />
  }
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            {t('profile:groupCaption', 'Group')}: {title}
          </div>
        </div>
        <div className='p-4 bg-purple-300 border border-purple-300 rounded w-1-2'>
            <p className='text-purple-900'>
                <strong>{t('profile:owner', "Owner")}:</strong>
            </p>
            <div className="px-2">
              <p>{t('profile:lastname', "Last Name")} {ownerLastName}</p>
              <p>{t('profile:firstname', "First Name")} {ownerFirstName}</p>
              <p>{t('profile:email', "Email")} {ownerEmail}</p>
            </div>
            <p className='mt-3 text-purple-900'>
                <strong>{t('profile:sharelink', "Share link for teachers")}:</strong>
            </p>
            <p className="px-2">
                <a href={BASE_API_URL + "/accept-group-link/" + shareLink}>{BASE_API_URL + "/accept-group-link/" + shareLink}</a>
            </p>

            <div className="flex items-center justify-between mt-2">
                <MultiSelect
                  className='relative block w-2/3 px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
                  options={userOptions}
                  value={selected}
                  onChange={setSelected}
                  labelledBy="Select"
                  filterOptions={filterOptions}
                  debounceDuration={500}
                />
                <button
                  className='relative flex justify-center px-2 py-1 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                  onClick={() => onSubmit()}
                >
                  {t('group:addStudent', 'Add Student')}
                </button>
            </div>

            { teachers.length !== 0 ? (
                <div className="mt-2">
                    <p className='mt-3 text-purple-900'>
                        <strong>{t('profile:teachers', "Teachers")}</strong>
                    </p>
                    <ul className="px-2">
                        {teachers.map((teacher, idx) => {
                            return <li key={idx}>{teacher.LastName} {teacher.FirstName} ({teacher.Email})</li>
                        })}
                    </ul>
                </div>
            ) : 
                null
            }

            <div className="mt-2">
                <p className='mt-3 text-purple-900'>
                    <strong>{t('profile:students', "Students")}</strong>
                </p>
                <ul className="px-2">
                    {members.map((member, idx) => {
                        return <li key={idx}>{member.LastName} {member.FirstName} ({member.Email})</li>
                    })}

                    {members.length === 0 ? <p>{t('profile:nogroups', "No students signed up for the group")}</p> : null }
                </ul>
            </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Group