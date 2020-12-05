import Footer from '../footer';
import Header from '../header';
import React from 'react';
import axios from 'axios';
import { Redirect } from 'react-router-dom';
import { connect } from "react-redux";
import { gradeHomework } from "../../actions/homework";
import authHeader from "../../services/auth-header";

class TeacherHwPage extends React.Component{
    
    constructor(props) {
        super(props);
        this.onChangeGrade = this.onChangeGrade.bind(this);
        this.onChangeComments = this.onChangeComments.bind(this);
        this.handleGrade = this.handleGrade.bind(this);
        this.ckEditorRemoveTags = this.ckEditorRemoveTags.bind(this);
        this.state = {
            id: null,
            grade: "",
            comments: "",
            title: "",
            description: "",
            files: [],
            closeDate: new Date(),
            homeworks: [],
            hwPageID: null
        };          
    }
    onChangeGrade(e){
        this.setState({
            grade: e.target.value,
            comments: "",
            isGraded: false
        });
    }

    onChangeComments(e){
        this.setState({
            comments: e.target.value
        }); 
    }
    ckEditorRemoveTags (data) {     
        const editedData = data.replace('<p>', '').replace('</p>', '') 
        return editedData;
    }
    handleGrade(id, index, e){
        console.log(index);
        console.log(this.state.homeworks[index].homework_page_id)
        e.preventDefault();
        const { dispatch} = this.props;
        dispatch(gradeHomework(this.state.homeworks[index].id, this.state.homeworks[index].student_fullname, this.state.homeworks[index].content, this.state.homeworks[index].submitted_at, this.state.homeworks[index].grade, this.state.homeworks[index].comments, this.state.hwPageID))
            .then(() => {
                this.setState({
                    isGraded: true,
                    successful: true
                });
            })
            .catch(() => {
                this.setState({
                    isGraded: false,
                    successful: false
                });
            });
    }

    componentDidMount () {
        const { randomStr } = this.props.match.params;
        axios.get(`/api/v1/homework-page/teacher/${randomStr}`, { headers: authHeader() })
            .then((response) => {
                if (response.data) {
                    this.setState({
                        hwPageID: response.data.id,
                        title: response.data.title,
                        description: response.data.content,
                        closeDate: response.data.closed_at,
                        homeworks: response.data.homeworks
                    })
                }
            })
      }
    render(){
        const data = this.ckEditorRemoveTags(this.state.description);
        const isEmptyDesc = this.state.description.trim() === "";
        const isEmptyFile = this.state.files.length === 0; 
        const submitted_hw = this.state.homeworks.length;
        const { isLoggedIn, role } = this.props;

        if (isLoggedIn && role === "student") {
            return <Redirect to="/"/>;
        }

        return(
            <div>
                <Header />
                <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                    <div className="flex justify-center flex-col items-center pb-4"></div>
                        <div className="text-purple-900 font-bold text-xl pt-4">Homework Submissions</div>
                    
                        <div className="w-3/4 border border-purple-300 rounded bg-purple-300 p-4">
                            <div className="flex flex-row"> 
                                <div className="w-3/4 flex flex-col">
                                    <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                        <strong>Homework Title:</strong> <span className="text-purple-900">{this.state.title}</span></h2> 
                                    {
                                        isEmptyDesc ? null :
                                        <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                            <strong>Description:</strong><br></br> <span className="text-purple-900">{data}</span></h2>
                                    }
                                    {
                                        isEmptyFile ? null :
                                        <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                            <strong>Attachments:</strong><br></br> <span className="text-purple-900">{this.state.files}</span></h2>
                                    } 
                                    <div className="flex flex-col items ml-2 ">
                                    {
                                        this.state.homeworks.map((homework, index) => (
                                            <div key={index} className="border border-purple-700 rounded flex flex-col mb-2">
                                                <div className="flex flex-col items  pb-2">
                                                    <p className="block tracking-wide text-gray-700 text-xs  px-4 pt-1">
                                                        <strong>Student's Name:</strong> <span className="text-purple-900">{homework.student_fullname}</span></p>
                                                    <p className="block tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1">
                                                        <strong>Submitted at:</strong> <span className="text-purple-900">{homework.submitted_at}</span></p>  
                                                    {
                                                        homework.content.trim() === "" ? null :
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Content:</strong><br /> <span className="text-purple-900">{this.ckEditorRemoveTags(homework.content)}</span></p>
                                                    }
                                                    {
                                                       /* homework.attachments.length === 0 ? null :
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Attachments:</strong><span className="text-purple-900">{homework.attachments}</span></p>*/
                                                    }       
                                                    {
                                                        homework.grade.trim() === "" ?
                                                        <form onSubmit={this.handleGrade.bind(this, homework.id, index)} className="flex flex-col">
                                                            <div className="flex flex-row mb-2">
                                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                                    <strong>Grade: </strong></p>
                                                                <input  onChange={this.onChangeGrade} className=" w-10 text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                                            </div>
                                                            <div className="flex flex-row">
                                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                                    <strong>Comments: </strong></p>
                                                                <textarea  onChange={this.onChangeComments} className=" text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                                                <br />
                                                            </div>
                                                            <div>
                                                                <button type="submit" className="mb-2 ml-4 mt-2
                                                                relative flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                                                    Grade Homework
                                                                </button>
                                                            </div>
                                                        </form> :
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Grade:</strong><span className="text-purple-900">{homework.grade}</span>
                                                            <strong>Comments:</strong><span className="text-purple-900">{homework.comments}</span>
                                                        </p> 
                                                    }
                                                </div>
                                            </div>
                                        ))
                                    }
                                    </div>
                                   
                                </div>
                                <div className="flex flex-col ">
                                    <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                       <strong><span className="text-purple-800 font-bold">{submitted_hw}</span> Homeworks Submitted</strong> </p>
                                    <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                        <strong></strong> </p>
                                    {/*<button type="submit" className="mb-2 ml-4 mt-2
                                                relative flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                        Close Homework
                                    </button>*/}
                                </div>
                            </div>
                        </div>
                </div>   
                <Footer />
            </div>
        )
    }
    }

function mapStateToProps(state) {
    const { isLoggedIn, role } = state.auth;
    return {
        isLoggedIn,
        role
    };
}

export default connect(mapStateToProps)(TeacherHwPage);