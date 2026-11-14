import { StyleSheet } from "react-native";

export default StyleSheet.create({
    safeArea: {
        backgroundColor: "#0A5133",
        flex: 1,
        justifyContent: 'center',
        alignItems: "center"
    },
    inputArea: {
        flex: 1,
        justifyContent: 'center',
        alignItems: "center",
        width: '100%',
        paddingHorizontal: 30,
    },
    singInPut: {
        width: '100%',
        height: 60,
        backgroundColor: '#A7F3D0', 
        flexDirection: 'row',
        borderRadius: 20, 
        paddingLeft: 15,
        alignItems: 'center',
        marginTop: 15,
        borderWidth: 2, 
        borderColor: '#1B5E20',
    },
    textInput: {
        flex: 1,
        fontSize: 16,
        marginLeft: 5,
        color: '#0A5133',
    },
    custtomButton: {
        width: '100%', 
        height: 60,
        backgroundColor: '#268596', 
        borderRadius: 20,
        justifyContent: 'center',
        alignItems: 'center',
        borderColor: '#1B5E20',
        borderWidth: 2, 
        marginTop: 30, 
        marginBottom: 20,
    },
    customButtontext: {
        fontSize: 18,
        color: '#FFFFFF',
        fontWeight: 'bold',
    },
    singMessageButto: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginVertical: 25,
    },
    singMessageButtonText: {
        fontSize: 15,
        color: '#E0E0E0', 
    },
    singMessageButtonTextBold: {
        fontSize: 15,
        fontWeight: 'bold',
        marginLeft: 5,
        color: "#A7F3D0",
    },
    loadingOverlay: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: 'rgba(0,0,0, 0.7)',
        zIndex: 1,
    },
});