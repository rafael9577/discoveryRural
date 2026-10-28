import { StyleSheet } from "react-native";

export default StyleSheet.create({
    SafeAreaView: {
        backgroundColor: "#0A5133",
        flex: 1,
    },
    busca: {
        alignSelf: "center",
        borderRadius: 20,
        width: "95%",
    },
    buscaInput: {
        marginTop: 20,
        backgroundColor: "gray",
        borderRadius: 20,
        borderWidth: 1,
        paddingHorizontal: 15,
        height: 45,
    },
    cardContainer: {
        backgroundColor: "#2E7D32",
        flexDirection: "row",
        alignItems: "center",
        padding: 20,
        borderWidth: 2,
        width: "95%",
        alignSelf: "center",
        marginTop: 20,
        borderRadius: 20,
    },
    cardTouchable: {
        flex: 1,
        flexDirection: "row",
        alignItems: "center",
    },
    cardImage: {
        backgroundColor: "pink",
        width: 100,
        height: 100,
        borderRadius: 20,
    },
    cardTextArea: {
        flex: 1,
        justifyContent: "center",
    },
    cardTitle: {
        marginLeft: 10,
        fontWeight: "bold",
        fontSize: 16,
        color: "#FFFFFF",
    },
    cardDescription: {
        marginLeft: 10,
        marginTop: 4,
        color: "#FFFFFF",
    },

    centerContainer: {
        alignItems: "center",
        justifyContent: "center",
        padding: 30,
    },
    loadingText: {
        color: "#FFFFFF",
        marginTop: 10,
    },
    errorText: {
        color: "#FFFFFF",
        textAlign: "center",
        marginBottom: 12,
    },
    emptyText: {
        color: "#FFFFFF",
        textAlign: "center",
    },
    retryButton: {
        backgroundColor: "#268596",
        paddingVertical: 10,
        paddingHorizontal: 20,
        borderRadius: 20,
    },
    retryButtonText: {
        color: "#FFFFFF",
        fontWeight: "bold",
    },
});