import { StyleSheet } from "react-native";

export default StyleSheet.create({
    safeArea: {
        backgroundColor: "#0A5133",
        flex: 1,
        justifyContent: 'center',
        alignItems: "center"
    },
    inputArea: {
        flex:1,
        justifyContent: 'center',
        alignItems: "center",
        width: '100%',
        padding: 40,
    },
    singInPut: {
        width: '100%',
        height: 60,
        backgroundColor: '#A7F3D0',
        flexDirection: 'row',
        borderRadius: 30,
        paddingLeft: 15,
        alignItems: 'center',
        marginTop: 10,
    },
    custtomButton: {
        height: 60,
        backgroundColor: '#268596',
        borderRadius: 30,
        justifyContent: 'center',
        alignItems: 'center',
        borderColor: 'black',
        borderWidth:0.6,
    },
    customButtontext: {
        fontSize: 18,
        color: '#ffffff',
    },
    singMessageButto: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginVertical: 20,
    },
    singMessageButtonTextBold: {
        fontSize: 16,
        color: 'black',
        fontWeight: 'bold',
        marginLeft: 5,
    },
    textInput: {
        flex: 1,
        fontSize: 16,
        marginLeft: 10,
    },
    forgotPassowrd: {
        marginRight:20,
        marginBottom: 15,
        marginTop: 10,
        alignSelf: "flex-end"
    },
    loadingOverlay: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: 'rgba(0,0,0, 0.5)',
        zIndex: 1,
    },
})