import { StyleSheet, Dimensions } from "react-native";

const { width, height } = Dimensions.get("window");

export default StyleSheet.create({
    SafeAreaView: {
        backgroundColor: "#0A5133",
        flex: 1,
    },
    title: {
        color: "#FFFFFF",
        fontSize: 22,
        fontWeight: "bold",
        marginBottom: 20,
        marginLeft: 5,
    },
    formContainer: {
        flex: 1,
        width: "95%",
        alignSelf: "center",
        marginTop: 20,
    },
    input: {
        backgroundColor: "#2E7D32",
        borderRadius: 20,
        borderWidth: 1,
        borderColor: "#1B5E20",
        paddingHorizontal: 15,
        height: 50,
        marginBottom: 15,
        color: "#FFFFFF",
    },
    textArea: {
        height: 100,
        textAlignVertical: "top",
        paddingTop: 15,
    },
    locationCard: {
        backgroundColor: "#2E7D32",
        padding: 20,
        borderRadius: 20,
        borderWidth: 2,
        borderColor: "#1B5E20",
        marginTop: 10,
        alignItems: "center",
    },
    locationTitle: {
        color: "#FFFFFF",
        fontWeight: "bold",
        fontSize: 16,
        marginBottom: 5,
    },
    locationCoords: {
        color: "#E0E0E0",
        fontSize: 14,
        marginBottom: 15,
        textAlign: "center",
    },
    mapButton: {
        backgroundColor: "#268596",
        paddingVertical: 12,
        paddingHorizontal: 25,
        borderRadius: 20,
        width: "100%",
    },
    mapButtonText: {
        color: "#FFFFFF",
        fontWeight: "bold",
        textAlign: "center",
        fontSize: 16,
    },
    submitButton: {
        backgroundColor: "#1B5E20",
        margin: 20,
        paddingVertical: 15,
        borderRadius: 20,
        borderWidth: 1,
        borderColor: "#FFFFFF",
    },
    submitButtonText: {
        color: "#FFFFFF",
        fontWeight: "bold",
        textAlign: "center",
        fontSize: 18,
    },
    
    centerContainer: {
        position: 'absolute',
        top: 0, left: 0, right: 0, bottom: 0,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: 'rgba(0,0,0, 0.7)',
        zIndex: 10,
    },
    loadingText: {
        color: "#FFFFFF",
        marginTop: 10,
        fontWeight: "bold",
    },

    mapContainer: {
        flex: 1,
    },
    map: {
        width: width,
        height: height,
    },
    fixedMarker: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        marginLeft: -10,
        marginTop: -30,
        alignItems: 'center',
    },
    markerDot: {
        width: 20,
        height: 20,
        backgroundColor: '#FF3B30',
        borderRadius: 10,
        borderWidth: 2,
        borderColor: '#FFFFFF',
    },
    markerStick: {
        width: 3,
        height: 20,
        backgroundColor: '#FF3B30',
    },
    mapFooter: {
        position: 'absolute',
        bottom: 40,
        left: 20,
        right: 20,
        flexDirection: 'row',
        justifyContent: 'space-between',
    },
    cancelMapButton: {
        backgroundColor: 'rgba(0,0,0,0.6)',
        paddingVertical: 15,
        borderRadius: 20,
        flex: 0.45,
    },
    cancelMapButtonText: {
        color: "#FFFFFF",
        fontWeight: "bold",
        textAlign: "center",
    },
    confirmMapButton: {
        backgroundColor: "#268596",
        paddingVertical: 15,
        borderRadius: 20,
        flex: 0.5,
    },
    confirmMapButtonText: {
        color: "#FFFFFF",
        fontWeight: "bold",
        textAlign: "center",
        fontSize: 16,
    },
});